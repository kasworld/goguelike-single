


# go get -d github.com/kasworld/h4o

################################################################################
cd lib
echo "genlog -leveldatafile ./g2log/g2log.data -packagename g2log "
genlog -leveldatafile ./g2log/g2log.data -packagename g2log 
cd ..


################################################################################
# generate enum
echo "generate enums"
genenum -typename=TurnAction          -packagename=turnaction          -basedir=enum -vectortype=int
genenum -typename=ReturnCode          -packagename=returncode          -basedir=enum -vectortype=int
genenum -typename=FlowType            -packagename=flowtype            -basedir=enum 
genenum -typename=AchieveType         -packagename=achievetype         -basedir=enum -vectortype=float64
genenum -typename=AIPlan              -packagename=aiplan              -basedir=enum -vectortype=int
genenum -typename=ActiveObjType       -packagename=aotype              -basedir=enum -vectortype=int
genenum -typename=CarryingObjectType  -packagename=carryingobjecttype  -basedir=enum -vectortype=int
genenum -typename=ClientControlType   -packagename=clientcontroltype   -basedir=enum 
genenum -typename=Condition           -packagename=condition           -basedir=enum -flagtype=uint16 -vectortype=int
genenum -typename=DangerType          -packagename=dangertype          -basedir=enum -vectortype=int
genenum -typename=DecayType           -packagename=decaytype           -basedir=enum
genenum -typename=EquipSlotType       -packagename=equipslottype       -basedir=enum -vectortype=int
genenum -typename=FactionType         -packagename=factiontype         -basedir=enum -vectortype=int
genenum -typename=FieldObjActType     -packagename=fieldobjacttype     -basedir=enum -vectortype=int
genenum -typename=FieldObjDisplayType -packagename=fieldobjdisplaytype -basedir=enum
genenum -typename=PotionType          -packagename=potiontype          -basedir=enum -vectortype=int
genenum -typename=ResourceType        -packagename=resourcetype        -basedir=enum -vectortype=int
genenum -typename=RespawnType         -packagename=respawntype         -basedir=enum 
genenum -typename=ScrollType          -packagename=scrolltype          -basedir=enum -vectortype=int
genenum -typename=StatusOpType        -packagename=statusoptype        -basedir=enum
genenum -typename=TerrainCmd          -packagename=terraincmd          -basedir=enum -vectortype=int
genenum -typename=Tile                -packagename=tile                -basedir=enum -flagtype=uint16 -vectortype=int
genenum -typename=TowerAchieve        -packagename=towerachieve        -basedir=enum -vectortype=float64
genenum -typename=TurnResultType      -packagename=turnresulttype      -basedir=enum
genenum -typename=Way9Type            -packagename=way9type            -basedir=enum 

goimports -w enum

$Data_VERSION=makesha256sum config/gameconst/*.go config/gamedata/*.go enum/*.enum
echo "Data Version: ${Data_VERSION}"
mkdir -ErrorAction SilentlyContinue config/dataversion
echo "package dataversion
const DataVersion = `"${Data_VERSION}`" 
" > config/dataversion/dataversion_gen.go 


################################################################################
$DATESTR=Get-Date -UFormat '+%Y-%m-%dT%H:%M:%S%Z:00'
$GITSTR=git rev-parse HEAD
################################################################################
# build bin

$BIN_DIR="bin"
$SRC_DIR="rundriver"

mkdir -ErrorAction SilentlyContinue "${BIN_DIR}"

# build bin here
$BUILD_VER="${DATESTR}_${GITSTR}_release_windows"
echo "Build Version: ${BUILD_VER}"
echo ${BUILD_VER} > ${BIN_DIR}/BUILD_windows
go build -o "${BIN_DIR}\goguelike-single.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\goguelike-single.go"
# go build -o "${BIN_DIR}\glclient.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\glclient.go"


# $BUILD_VER="${DATESTR}_${GITSTR}_release_linux"
# echo "Build Version: ${BUILD_VER}"
# echo ${BUILD_VER} > ${BIN_DIR}/BUILD_linux
# $env:GOOS="linux" 
# go build -o "${BIN_DIR}\goguelike-single" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\goguelike-single.go"
# go build -o "${BIN_DIR}\glclient" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\glclient.go"
# $env:GOOS=""

# $BUILD_VER="${DATESTR}_${GITSTR}_release_wasm"
# cd rundriver
# ./genwasmclient.ps1 ${BUILD_VER}
# cd ..

echo "cp -r rundriver/serverdata ${BIN_DIR}"
Copy-Item -Force -r rundriver/serverdata ${BIN_DIR}
echo "cp -r rundriver/clientdata ${BIN_DIR}"
Copy-Item -Force -r rundriver/clientdata ${BIN_DIR}

