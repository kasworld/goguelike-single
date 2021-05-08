# Modified G3N for my taste

## change

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
