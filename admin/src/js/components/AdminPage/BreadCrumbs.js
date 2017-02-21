import React, {PropTypes} from 'react'
import { Link } from 'react-router'

class BreadCrumbs extends React.Component {

  render () {
    const { 
      pathname,
      breadcrumbs
    } = this.props

    const breadCrumbsLength = breadcrumbs.filter(l => { return !!l }).size

    return (
      <ul className="bread-crumbs">
        { 
          breadcrumbs
            .filter((l, i) => {
              return !!l && i < breadCrumbsLength - i
            })
            .map((l, i) => {
              return (
                <li
                  key={ 'breadcrumb-' + i }
                  className="bread-crumbs_level"
                >
                  <Link to={ l.get('url') }>
                    { l.get('name') }
                  </Link>
                </li>
              )
            })
        }
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