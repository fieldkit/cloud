import React, {PropTypes} from 'react'
import { Link } from 'react-router'
import I from 'immutable'


class NewConfirmationSection extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      addMemberValue: null,
      inputValues: {}
    }
  }

  render () {

    const { 
      currentProjectID,
      currentExpedition
    } = this.props

    return (
      <div id="new-inputs-section" className="section">
        <div className="section-header">
          <h1>Your connected expedition is ready</h1>
        </div>
        
        <p className="intro">
          Congratulations, you're done creating your first expedition! You can now dive deeper in the settings, or go straight see how your map looks like.
        </p>
       
        <div className="call-to-action"> 
        <Link to={'/admin/' + currentProjectID + '/' + currentExpedition.get('id') }>
          <div className="button">
            Go to the admin dashboard
          </div>
        </Link>
        or
        <a href={'https://' + currentProjectID + '.fieldkit.org/' + currentExpedition.get('id') }>
          <div className="button hero">
            Go to my map!
          </div>
        </a>
        </div>

      </div>
    )

  }
}

NewConfirmationSection.propTypes = {
}

export default NewConfirmationSection
