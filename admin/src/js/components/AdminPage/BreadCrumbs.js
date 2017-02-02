import React, {PropTypes} from 'react'

class BreadCrumbs extends React.Component {

  render () {
    const { 
      pathname,
      breadcrumbs
    } = this.props

    console.log('LOOOL', breadcrumbs)

    const levels = breadcrumbs
      .filter((l, i) => {
        return i < pathname.split('/').length - 2
      })
      .map((l, i) => {
        return (
          <li
            key={i}
            className="bread-crumbs_level"
          >
            { l }
          </li>
        )
      })

    return (
      <ul className="bread-crumbs">
        { levels }
      </ul>
    )

      /*
      <Route path="profile"/>

      <Route path="new-project"/>

      <Route path=":projectID">
        <Route path="new-expedition">
          <Route path="general-settings" component={NewGeneralSettingsContainer}/>
          <Route path="inputs" component={NewInputsContainer}/>
          <Route path="confirmation" component={NewConfirmationContainer}/>
        </Route>
        <Route path=":expeditionID">
          <Route path="dashboard" component={ExpeditionPageContainer}/>
          <Route path="general-settings" component={GeneralSettingsContainer}/>
          <Route path="inputs" component={InputsContainer}/>
        </Route>
      </Route>
    */

    return (
      <div></div>
    )
  }

}

BreadCrumbs.propTypes = {}

export default BreadCrumbs