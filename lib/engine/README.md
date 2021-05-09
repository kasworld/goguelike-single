# Modified G3N for my taste

## requirement and install

to auto fix import 
    
    goimports 
    go get golang.org/x/tools/cmd/goimports 

to modify eventtype enum (https://github.com/kasworld/genenum)

    go get github.com/kasworld/genenum
    

to modify loglevel (https://github.com/kasworld/log)

    go get github.com/kasworld/log
    run install.sh to make genlog


## change from G3N (https://github.com/g3n/engine)

remove wasm support 

rename interface I*** -> ***I 

    NodeI
    AppWindowI
    DispatcherI
    MaterialI
    LightI
    GraphicI
    GeometryI
    CameraI
    LayoutI
    PanelI
    MaterialI
    BuilderLayoutI
    SolverI
    BodyI
    EquationI
    ConstraintI
    ForceFieldI
    ShapeI
    ChannelI

split package, file 

    split interface file 
    core -> node, dispatcher, renderinfo, timemanager
    color -> colornames
    array -> array_f32, array_u32
    util -> framerater

del app singletone 

rename app to appbase, Application to AppBase

    use same package name with main struct name as possible 

merge event to one file and to enum by genenum

    evname string -> evname eventtype.EventType

replace logger 

rename window -> appwindow

## TODO 

split interface package more 
del no need interface