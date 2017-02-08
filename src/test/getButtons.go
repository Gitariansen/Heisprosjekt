//get Buttons
package test
import (
	"fmt"
  "time"
  "../driver"
)

/*Iterate over all floors, all types of buttons (through the button_channel_matrix) and add to channel*/
func Get_Button_Press () /*<-chan bool*/{
  //ret := make(chan bool, 5)
  go func() {
    for {
      for floor := 0; floor < driver.N_FLOORS; floor++{
        for button := 0; button < driver.N_BUTTONS; button ++{
					fmt.Println("New iteration", floor, button)
					fmt.Println(driver.Elev_get_button_signal(button, floor))
        }
      }
			time.Sleep(time.Millisecond)
    }
  }()
	//return ret
}
