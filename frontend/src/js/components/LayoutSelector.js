
import React, {PropTypes} from 'react'

const LayoutSelector = ({mode, onLayoutChange}) => {
  var types = ['rows','grid']

  var buttons = types.map(function(s,i){

    var className = 'layoutButton ' + (s === mode ?'active':'inactive')

    return (
      <li className={className} key={i} onClick={()=>onLayoutChange(s)}>
        <img width="16" height="16"/>
      </li>
    )
  })
  return (
    <div className="selector">
      <div className="column">
        <ul className="buttonRow">
          {buttons}
        </ul>
      </div>
    </div>
  )
}

LayoutSelector.propTypes = {
  onLayoutChange: PropTypes.func.isRequired,
  mode: PropTypes.string.isRequired
}

export default LayoutSelector
