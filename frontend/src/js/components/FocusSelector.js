
import React, {PropTypes} from 'react'

const FocusSelector = ({mainFocus, secondaryFocus, onMainFocusChange, onSecondaryFocusChange}) => {

  var toggleDropdown = (i) => {
    document.getElementById("FocusSelector"+i+"Options").classList.toggle("show");
  }

  return (
    <div className="focusSelector controlSelector">
      <p>Focus on:</p>
      <div className="dropdown">
        <button onClick={()=>toggleDropdown(1)} className="dropbtn">{mainFocus}</button>
        <div id="FocusSelector1Options" className="dropdown-content">
          <a href="#" onClick={()=>onMainFocusChange("explorers")}>Explorers</a>
          <a href="#" onClick={()=>onMainFocusChange("sensors")}>Sensors</a>
          <a href="#" onClick={()=>onMainFocusChange("animals")}>Animals</a>
        </div>
      </div>
      <div className="dropdown">
        <button onClick={()=>toggleDropdown(2)} className="dropbtn">{secondaryFocus}</button>
        <div id="FocusSelector2Options" className="dropdown-content">
          <a href="#" onClick={()=>onSecondaryFocusChange("steve")}>Steve</a>
          <a href="#" onClick={()=>onSecondaryFocusChange("jer")}>Jer</a>
          <a href="#" onClick={()=>onSecondaryFocusChange("shah")}>Shah</a>
          <a href="#" onClick={()=>onSecondaryFocusChange("chris")}>Chris</a>
        </div>
      </div>
    </div>
  )
}

FocusSelector.propTypes = {
  mainFocus: PropTypes.string.isRequired,
  secondaryFocus: PropTypes.string.isRequired,
  onMainFocusChange: PropTypes.func.isRequired,
  onSecondaryFocusChange: PropTypes.func.isRequired
}

export default FocusSelector