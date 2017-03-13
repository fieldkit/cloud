import React, {PropTypes} from 'react'
import {findDOMNode} from 'react-dom'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
import ContentEditable from 'react-contenteditable'
import I from 'immutable'
import Dropdown from 'react-dropdown'
import Select from 'react-select'
import { Base64 } from 'js-base64'
import { protocol, hostname } from '../../../constants/APIBaseURL'

import iconRemoveSmall from '../../../../img/icon-remove-small.png'


class ThemeSection extends React.Component {  
  render () {
    const { 
      projectID,
      expedition,
      errors,
      setExpeditionProperty
    } = this.props

    return (
      <div id="new-inputs-section" className="section">
        <div className="section-header">
          <h1>Theme Settings</h1>
        </div>
        <p className="intro">
          This page will be used to set up the expedition's theme and visual identity.
        </p>

        <div className="status">
        </div>

        {
          !!errors &&
          <p className="errors">
            We found one or multiple errors. Please check the information above or try again later.
          </p>
        }
      </div>
    )
  }
}

ThemeSection.propTypes = {
}

export default ThemeSection
