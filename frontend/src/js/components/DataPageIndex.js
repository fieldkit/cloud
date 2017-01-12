
import React, {PropTypes} from 'react'

const DataPageIndex = ({children}) => (
  <div id="APIIndex">
    {children}
  </div>
)

DataPageIndex.propTypes = {
  children: PropTypes.node.isRequired
}

export default DataPageIndex