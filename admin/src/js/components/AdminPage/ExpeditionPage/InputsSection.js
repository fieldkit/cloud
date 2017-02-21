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


class NewInputsSection extends React.Component {
  render () {
    const { 
      projectID,
      expedition,
      documentTypes,
      selectedDocumentType,
      setExpeditionProperty,
      fetchSuggestedDocumentTypes,
      addDocumentType,
      removeDocumentType,
      submitInputs
    } = this.props

    return (
      <div id="new-inputs-section" className="section">
        <div className="section-header">
          <h1>Set up your devices (2/2)</h1>
        </div>
        <p className="intro">
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien. Praesent bibendum leo vitae felis pellentesque, sit amet mattis nisi mattis.
        </p>

        <h5>Sensor devices</h5>
        <p>
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. 
        </p>
        <table className="objects-list">
          {
            !!documentTypes && !!documentTypes.filter(d => {return d.get('type') === 'sensor'}).size &&
            <tbody>
              <td className="name">Name</td>
              <td className="description">Description</td>
              <td className="inputs">Setup instructions</td>
              <td className="remove"></td>
            </tbody>
          }
          { 
            documentTypes
              .filter(d => {return d.get('type') === 'sensor'}) 
              .map((d, i) => {
                return (
                  <tbody key={i}>
                    <td className="name">{d.get('name')}</td>
                    <td className="description">{d.get('description')}</td>
                    <td className="inputs">{ (() => {
                        switch(d.get('setupType')) {
                          case 'token' : {
                            return (
                              <ul>
                                <li>
                                  1. Download this configuration file on a SD card:
                                </li>
                                <li>
                                  <a 
                                    href={'data:application/octet-stream;charset=utf-16le;base64,' + Base64.encode(
                                      protocol + hostname + '/api/input/' + d.get('token') + '/fieldkit/rockblock?token=' + expedition.get('token')
                                    )}
                                    download="config.txt"
                                  >
                                    <div className="button">Download config file</div>
                                  </a>
                                </li>
                                <li>
                                  2. Initialize your device. You're done!
                                </li>
                              </ul>
                            )
                          }
                        }
                      })()
                    }</td>
                    <td 
                      className="remove"
                      onClick={() => {
                        removeDocumentType(d.get('id'))
                      }}
                    >  
                      <img src={'/' + iconRemoveSmall}/>
                    </td>
                  </tbody>
                )
              })
          }
          <tbody>
            <td className="add-object" colSpan="2" width="60%">
              <div className="add-object-container">
                <Select.Async
                  name="add-documentType"
                  loadOptions={(input, callback) =>
                    fetchSuggestedDocumentTypes(input, 'sensor', callback)
                  }
                  value={ expedition.getIn(['selectedDocumentType', 'sensor']) }
                  onChange={(val) => {
                    setExpeditionProperty(['selectedDocumentType', 'sensor'], val.value)
                  }}
                  clearable={false}
                />
                <div
                  className={ "button" + (!!expedition.getIn(['selectedDocumentType', 'sensor']) ? '' : ' disabled') }
                  onClick={() => {
                    if (!!expedition.getIn(['selectedDocumentType', 'sensor'])) {
                      addDocumentType(expedition.getIn(['selectedDocumentType', 'sensor']), 'sensor')
                    }
                  }}
                >
                  Add Document Type
                </div>
              </div>
            </td>
            <td className="add-object-label" colSpan="3" width="50%">
              Search by document name (e.g tweet, geolocation...)
            </td>
          </tbody>
        </table>
        {
          !!errors && !!errors.get && !!errors.get('name') &&
          errors.get('name').map((error, i) => {
            return (
              <p 
                key={'errors-name-' + i}
                className="errors"
              >
                {error}
              </p>
            )
          })
        }

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

NewInputsSection.propTypes = {
}

export default NewInputsSection
