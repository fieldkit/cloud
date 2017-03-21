
import React from 'react'

import iconMinus from '../../../../img/icon-minus.png'
import iconPlus from '../../../../img/icon-plus.png'

class ZoomSelector extends React.Component {
  render () {
    const {
      zoom,
      selectZoom
    } = this.props

    return (
      <ul className="control-panel_zoom-selector">
        <li
          className="control-panel_zoom-selector_button"
          onClick={ () => selectZoom(zoom - 1) }
        >
          <img
            src={ '/' + (iconMinus) }
          />
        </li>
        <li
          className="control-panel_zoom-selector_button"
          onClick={ () => selectZoom(zoom + 1) }
        >
          <img
            src={ '/' + (iconPlus) }
          />
        </li>
      </ul>
    )
  }
}

export default ZoomSelector