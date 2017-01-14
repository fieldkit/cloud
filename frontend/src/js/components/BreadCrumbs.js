import React, {PropTypes} from 'react'

class BreadCrumbs extends React.Component {

  render () {
    const { pathname } = this.props

    const path = pathname.split('/')
      .filter((p, i) => {
        return i > 0
      })
      .map(p => {
        
      })

    return (
      <p id="breadcrumbs">
        test1 / test2 / test3
      </p>
    )
  }

}

BreadCrumbs.propTypes = {}

export default BreadCrumbs