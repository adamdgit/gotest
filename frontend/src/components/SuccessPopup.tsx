import { Accessor, createEffect, Setter } from "solid-js";

export default function SuccessPopup({ message, showPopup, setShowPopup }: { 
  message: Accessor<string>,
  showPopup: Accessor<boolean>,
  setShowPopup: Setter<boolean>
}) {

    createEffect(() => {
      if (showPopup()) {
        const timer = setTimeout(() => {
          setShowPopup(false);
        },2500)
        
        return () => clearTimeout(timer)
      }
    })

  return (
    <div class="popup" style={showPopup() ? {"opacity": "1"} : {"opacity": "0"}}>
        {message()}
    </div>
  )
}
