package main


import (
    "fmt"
    "strings"
    "./network"
)


func main() {
    fmt.Println("Hella")
    IP := network.get_local_IP()
    fmt.Println(IP)
}
