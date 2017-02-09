//get Buttons
package test
import (
	"fmt"
  "time"
  "driver"
)

/*Iterate over all floors, all types of buttons (through the button_channel_matrix) and add to channel*/
func Get_Button_Press () /*<-chan bool*/{
  //ret := make(chan bool, 5)
  go func() {
    for {
      for floor := 0; floor < driver.N_FLOORS; floor++{
        for button := 0; button < driver.N_BUTTONS; button ++{
					if (driver.Elev_get_button_signal(button, floor) == 1){
						fmt.Println("Button Pressed in floor: ", floor + 1)
						for driver.Elev_get_button_signal(button, floor) == 1 {
							
						}
					}
        }
      }
			time.Sleep(100*time.Millisecond)
    }
  }()
	//return ret
}
