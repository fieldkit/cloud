
import React from 'react'

class ExpeditionPanel extends React.Component {
  render () {
    const {
      active,
      closeExpeditionPanel
    } = this.props

    return (
      <div
        className={`expedition-panel ${ active ? 'active' : '' }`}
        onClick={ closeExpeditionPanel }
      >
        <div
          className="content"
          onClick={ (e) => { e.stopPropagation() }}
        >
          This is the expedition panel...
        </div>
      </div>
    )
  }
}

export default ExpeditionPanel