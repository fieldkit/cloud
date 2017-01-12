import React, {PropTypes} from 'react'
import {findDOMNode} from 'react-dom'
import { Link } from 'react-router'
import autobind from 'autobind-decorator'
import ContentEditable from 'react-contenteditable'
import I from 'Immutable'
import Dropdown from 'react-dropdown'
import Select from 'react-select'

class NewInputsSection extends React.Component {
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
      currentExpedition,
      documentTypes,
      selectedDocumentType,
      setExpeditionProperty,
      fetchSuggestedDocumentTypes,
      addDocumentType,
      removeDocumentType
    } = this.props

    return (
      <div id="new-inputs-section" className="section">
        <div className="section-header">
          <h1>Configure document types</h1>
          <p>process breadcrumbs</p>
        </div>
        <p className="intro">
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. Mauris pretium, nunc non lacinia finibus, dui lectus molestie nulla, quis ultricies libero orci a sapien. Praesent bibendum leo vitae felis pellentesque, sit amet mattis nisi mattis.
        </p>

        <h5>Team recorded documents</h5>
        <p>
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. 
        </p>
        <table className="objects-list">
          {
            !!documentTypes && !!documentTypes.filter(d => {return d.get('type') === 'member'}).size &&
            <tbody>
              <td className="name">Name</td>
              <td className="description">Description</td>
              <td className="inputs">Inputs</td>
              <td className="remove"></td>
            </tbody>
          }
          { 
            documentTypes
              .filter(d => {return d.get('type') === 'member'}) 
              .map((d, i) => {
                return (
                  <tbody key={i}>
                    <td className="name">{d.get('name')}</td>
                    <td className="description">{d.get('description')}</td>
                    <td className="inputs">{
                      d.get('inputs')
                        .map((input, j) => {
                          return (
                            <span key={j}>{input}</span>
                          )
                        })
                    }</td>
                    <td 
                      className="remove"
                      onClick={() => {
                        removeDocumentType(d.get('id'))
                      }}
                    >  
                      <img src="/src/img/icon-remove-small.png"/>
                    </td>
                  </tbody>
                )
              })
          }
          <tbody>
            <td className="add-object" colSpan="3" width="50%">
              <div className="add-object-container">
                <Select.Async
                  name="add-documentType"
                  loadOptions={(input, callback) =>
                    fetchSuggestedDocumentTypes(input, 'member', callback)
                  }
                  value={ currentExpedition.getIn(['selectedDocumentType', 'member']) }
                  onChange={(val) => {
                    setExpeditionProperty(['selectedDocumentType', 'member'], val.value)
                  }}
                  clearable={false}
                />
                <div
                  className={ "button" + (!!currentExpedition.getIn(['selectedDocumentType', 'member']) ? '' : ' disabled') }
                  onClick={() => {
                    if (!!currentExpedition.getIn(['selectedDocumentType', 'member'])) {
                      addDocumentType(currentExpedition.getIn(['selectedDocumentType', 'member']), 'member')
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

        <h5>Deployed sensor documents</h5>
        <p>
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. 
        </p>
        <table className="objects-list">
          {
            !!documentTypes && !!documentTypes.filter(d => {return d.get('type') === 'sensor'}).size &&
            <tbody>
              <td className="name">Name</td>
              <td className="description">Description</td>
              <td className="inputs">Inputs</td>
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
                    <td className="inputs">{
                      d.get('inputs')
                        .map((input, j) => {
                          return (
                            <span key={j}>{input}</span>
                          )
                        })
                    }</td>
                    <td 
                      className="remove"
                      onClick={() => {
                        removeDocumentType(d.get('id'))
                      }}
                    >  
                      <img src="/src/img/icon-remove-small.png"/>
                    </td>
                  </tbody>
                )
              })
          }
          <tbody>
            <td className="add-object" colSpan="3" width="50%">
              <div className="add-object-container">
                <Select.Async
                  name="add-documentType"
                  loadOptions={(input, callback) =>
                    fetchSuggestedDocumentTypes(input, 'sensor', callback)
                  }
                  value={ currentExpedition.getIn(['selectedDocumentType', 'sensor']) }
                  onChange={(val) => {
                    setExpeditionProperty(['selectedDocumentType', 'sensor'], val.value)
                  }}
                  clearable={false}
                />
                <div
                  className={ "button" + (!!currentExpedition.getIn(['selectedDocumentType', 'sensor']) ? '' : ' disabled') }
                  onClick={() => {
                    if (!!currentExpedition.getIn(['selectedDocumentType', 'sensor'])) {
                      addDocumentType(currentExpedition.getIn(['selectedDocumentType', 'sensor']), 'sensor')
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

        <h5>Social media documents</h5>
        <p>
          Etiam eu purus in urna volutpat ornare. Etiam pretium ante non egestas dapibus. 
        </p>
        <table className="objects-list">
          {
            !!documentTypes && !!documentTypes.filter(d => {return d.get('type') === 'social'}).size &&
            <tbody>
              <td className="name">Name</td>
              <td className="description">Description</td>
              <td className="inputs">Inputs</td>
              <td className="remove"></td>
            </tbody>
          }
          { 
            documentTypes
              .filter(d => {return d.get('type') === 'social'}) 
              .map((d, i) => {
                return (
                  <tbody key={i}>
                    <td className="name">{d.get('name')}</td>
                    <td className="description">{d.get('description')}</td>
                    <td className="inputs">{
                      d.get('inputs')
                        .map((input, j) => {
                          return (
                            <span key={j}>{input}</span>
                          )
                        })
                    }</td>
                    <td 
                      className="remove"
                      onClick={() => {
                        removeDocumentType(d.get('id'))
                      }}
                    >  
                      <img src="/src/img/icon-remove-small.png"/>
                    </td>
                  </tbody>
                )
              })
          }
          <tbody>
            <td className="add-object" colSpan="3" width="50%">
              <div className="add-object-container">
                <Select.Async
                  name="add-documentType"
                  loadOptions={(input, callback) =>
                    fetchSuggestedDocumentTypes(input, 'social', callback)
                  }
                  value={ currentExpedition.getIn(['selectedDocumentType', 'social']) }
                  onChange={(val) => {
                    setExpeditionProperty(['selectedDocumentType', 'social'], val.value)
                  }}
                  clearable={false}
                />
                <div
                  className={ "button" + (!!currentExpedition.getIn(['selectedDocumentType', 'social']) ? '' : ' disabled') }
                  onClick={() => {
                    if (!!currentExpedition.getIn(['selectedDocumentType', 'social'])) {
                      addDocumentType(currentExpedition.getIn(['selectedDocumentType', 'social']), 'social')
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

        <p className="status">
        </p>

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
        <Link to={'https://' + currentProjectID + '.fieldkit.org/' + currentExpedition.get('id') }>
          <div className="button hero">
            Go to my map!
          </div>
        </Link>
        </div>

      </div>
    )

  }
}

NewInputsSection.propTypes = {
}

export default NewInputsSection
