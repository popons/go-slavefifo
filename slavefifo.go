package slavefifo

import "fmt";
import "libusb";
import "time";


// vendor id for primary-usb
var Vid = 0x04b4;
// product id for primary-usb
var Pid = 0x8613;
// vendor id for secondary-usb
var SVid = 0x04b4;
// product id for secondary-usb
var SPid = 0x1004;

var usbdev *libusb.Device = nil;

type Error struct
{
    emsg string;
}

func (e *Error) String() string     // Stringer interface
{
    return e.emsg;
}

func found(vid,pid int) bool
{
    for _,info := range libusb.Enum()
    {
        if info.Vid==vid && info.Pid==pid
        {
            return true;
        }
    }
    return false;
}

func open(vid,pid int , primary bool) *Error
{
    var r int;
    var e *Error;

    usbdev = libusb.Open(vid,pid);

    if usbdev == nil { return &Error{fmt.Sprintf("VID:%04X PID%04X can not open or not found. ",vid,pid)};}

    if !primary
    {
        r = usbdev.Configuration(1);
        if r!=0 { return &Error{fmt.Sprintf(" Configuration error %d [%s]",r,usbdev.LastError())}; }

        usbdev.Interface(0);
        if r!=0 { return &Error{fmt.Sprintf(" Interface error %d [%s]",r,usbdev.LastError())}; }
    }

    if primary
    {
        e = reset(usbdev,1);
        if e != nil { return e; }

        for i:=0;i<len(prog_data); i+= 4096
        {
            r := usbdev.ControlMsg(libusb.USB_TYPE_VENDOR,0xa0,i,0,prog_data[i:i+4096]);
            fmt.Printf("ControlMsg() return %d\n",r);
            if r!=4096 { return &Error{fmt.Sprintf(" Programming error %d [%s]",r,usbdev.LastError())}; }
        }

        e = reset(usbdev,0);
        if e != nil { return e; }
    }

    return nil;
}

func Open() *Error
{
    b,d := libusb.Init();
    if b==0 || d==0 { return &Error{ fmt.Sprintf("busN=%d devN=%d ",b,d) }; }

    if found(Vid,Pid)
    {
        e := open(Vid,Pid,true);
        if e != nil { return e; }

        e = Close();
        if e != nil { return e; }

    }
    for c:=0; c<100; c++
    {
        libusb.Init();
        if found(SVid,SPid)     // seach Seconday
        {
            Close();
            return open(SVid,SPid,false);
        }
        time.Sleep(1e8);    // 1e8ns 100ms
        fmt.Printf(".");
    }

    return &Error{fmt.Sprintf("Secondary Device VID%04X PID%04X not found.",SVid,SPid)};

}

func Close() *Error
{
    if usbdev != nil
    {
        r := usbdev.Close();
        usbdev = nil;
        if r != 0
        {
            return &Error{ fmt.Sprintf("libusb.Device.Close() return %d",r) };
        }
    }
    return nil;
}

func reset(device *libusb.Device,bit byte) *Error
{
    var dat = []byte{bit};
    r := device.ControlMsg(libusb.USB_TYPE_VENDOR,0xa0,0xe600,0,dat);
    fmt.Printf("usb_control_msg() return %d\n",r);
    if r!=1
    {
        return &Error{ fmt.Sprintf("reset to %d ControlMsg return %d in reset.",bit,r) }
    }
    return nil;
}

func Write(ep int,dat []uint32) *Error
{
    return nil;
}
func Read(ep int,dat []uint32) *Error
{
    return nil;
}
