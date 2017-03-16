import React, {PropTypes} from 'react'
import { Link } from 'react-router'

class ProfileSection extends React.Component {
  constructor (props) {
    super(props)
    this.state = {

    }
  }
  
  render () {

    return (
      <div id="profile-section" className="section">
        <h2>Profile Section</h2>
      </div>
    )
  }
}

ProfileSection.propTypes = {

}

export default ProfileSection