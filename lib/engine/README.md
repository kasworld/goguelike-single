# Modified G3N for my taste

## require to modify eventtype enum 

https://github.com/kasworld/genenum

## change from G3N (https://github.com/g3n/engine)

remove wasm support 

rename interface I*** -> ***I 

    NodeI
    WindowI
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

## TODO 

split interface package more 
del no need interface