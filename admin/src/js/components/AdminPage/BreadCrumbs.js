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
      <div></div>
    )
  }

}

BreadCrumbs.propTypes = {}

export default BreadCrumbs