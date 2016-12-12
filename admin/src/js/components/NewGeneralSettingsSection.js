import React, {PropTypes} from 'react'
import {findDOMNode} from 'react-dom'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
import ContentEditable from 'react-contenteditable'
import I from 'Immutable'
import Dropdown from 'react-dropdown'
import Select from 'react-select';

class NewGeneralSettingsSection extends React.Component {
  constructor (props) {
    super(props)
    this.state = {
      addMemberValue: null,
      inputValues: {}
    }
  }

  render () {

    const { } = this.props

    return (
      <p>NewGeneralSettingsSection</p>
    )

  }
}

NewGeneralSettingsSection.propTypes = {
}

export default NewGeneralSettingsSection