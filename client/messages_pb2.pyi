from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Empty(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class Positions(_message.Message):
    __slots__ = ("positions",)
    POSITIONS_FIELD_NUMBER: _ClassVar[int]
    positions: _containers.RepeatedCompositeFieldContainer[Position]
    def __init__(self, positions: _Optional[_Iterable[_Union[Position, _Mapping]]] = ...) -> None: ...

class Position(_message.Message):
    __slots__ = ("x", "y")
    X_FIELD_NUMBER: _ClassVar[int]
    Y_FIELD_NUMBER: _ClassVar[int]
    x: int
    y: int
    def __init__(self, x: _Optional[int] = ..., y: _Optional[int] = ...) -> None: ...

class Color(_message.Message):
    __slots__ = ("r", "g", "b")
    R_FIELD_NUMBER: _ClassVar[int]
    G_FIELD_NUMBER: _ClassVar[int]
    B_FIELD_NUMBER: _ClassVar[int]
    r: int
    g: int
    b: int
    def __init__(self, r: _Optional[int] = ..., g: _Optional[int] = ..., b: _Optional[int] = ...) -> None: ...

class Light(_message.Message):
    __slots__ = ("active", "inactive")
    ACTIVE_FIELD_NUMBER: _ClassVar[int]
    INACTIVE_FIELD_NUMBER: _ClassVar[int]
    active: Color
    inactive: Color
    def __init__(self, active: _Optional[_Union[Color, _Mapping]] = ..., inactive: _Optional[_Union[Color, _Mapping]] = ...) -> None: ...

class LightStatus(_message.Message):
    __slots__ = ("position", "status")
    POSITION_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    position: Position
    status: Light
    def __init__(self, position: _Optional[_Union[Position, _Mapping]] = ..., status: _Optional[_Union[Light, _Mapping]] = ...) -> None: ...

class LightsStatus(_message.Message):
    __slots__ = ("lights",)
    LIGHTS_FIELD_NUMBER: _ClassVar[int]
    lights: _containers.RepeatedCompositeFieldContainer[LightStatus]
    def __init__(self, lights: _Optional[_Iterable[_Union[LightStatus, _Mapping]]] = ...) -> None: ...

class SensorStatus(_message.Message):
    __slots__ = ("position", "status")
    POSITION_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    position: Position
    status: bool
    def __init__(self, position: _Optional[_Union[Position, _Mapping]] = ..., status: bool = ...) -> None: ...

class Brightness(_message.Message):
    __slots__ = ("position", "brightness")
    POSITION_FIELD_NUMBER: _ClassVar[int]
    BRIGHTNESS_FIELD_NUMBER: _ClassVar[int]
    position: Position
    brightness: int
    def __init__(self, position: _Optional[_Union[Position, _Mapping]] = ..., brightness: _Optional[int] = ...) -> None: ...

class BrightnessStatus(_message.Message):
    __slots__ = ("brightness",)
    BRIGHTNESS_FIELD_NUMBER: _ClassVar[int]
    brightness: _containers.RepeatedCompositeFieldContainer[Brightness]
    def __init__(self, brightness: _Optional[_Iterable[_Union[Brightness, _Mapping]]] = ...) -> None: ...
