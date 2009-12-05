package slavefifo

import "testing";
import "os";
import "fmt";

func TestSlaveFifo(t *testing.T)
{
    {
    e := Open();
    if e != nil
    {
        fmt.Printf("%v\n",e);
        os.Exit(1);
    }
    }
    {
    e := Close();
    if e != nil
    {
        fmt.Printf("%v\n",e);
        os.Exit(1);
    }
    }
}
