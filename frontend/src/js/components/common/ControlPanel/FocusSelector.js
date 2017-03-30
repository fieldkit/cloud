
import React from 'react'

class FocusSelector extends React.Component {

  shouldComponentUpdate (props) {
    return this.props.focus.get('type') !== props.focus.get('type') ||
      this.props.focus.get('id') !== props.focus.get('id')
  }

  render () {
    const {
      focus,
      selectFocusType
    } = this.props

    return (
      <div className="control-panel_focus-selector">
        <p>Map Focus:</p>
        <ul>
          <li
            className={ focus.get('type') === 'expedition' ? 'active' : '' }
            onClick={ () => selectFocusType('expedition') }
          >
            Expedition
          </li>
          <li
            className={ focus.get('type') === 'documents' ? 'active' : '' }
            onClick={ () => selectFocusType('documents') }
          >
            Documents
          </li>
          <li
            className={ focus.get('type') === 'manual' ? 'active' : '' }
            onClick={ () => selectFocusType('manual') }
          >
            Manual
          </li>
        </ul>
      </div>
    )
  }
}

export default FocusSelector