import React, {PropTypes} from 'react'

class Modal extends React.Component {

  render () {

    const { 
      modal,
      saveChangesAndResume,
      cancelAction
    } = this.props

    return (
      <div className="modal-container">
        <div className="modal">
          <p>
            Please save or cancel your changes to [] first.
          </p>
          <div className="actions">
            <div 
              className="button secondary"
              onClick={ cancelAction }
            >
              Cancel
            </div>
            <div 
              className="button primary"
              onClick={ saveChangesAndResume }
            >
              Save changes
            </div>
          </div>
        </div>
      </div>
    )
  }
}

Modal.propTypes = {
  saveChangesToTeam: PropTypes.func.isRequired,
  clearChangesToTeam: PropTypes.func.isRequired
}

export default Modal