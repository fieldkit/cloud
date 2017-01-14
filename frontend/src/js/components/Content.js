import React, { PropTypes } from 'react'


import { Route, IndexRoute } from 'react-router'


const Content = ({currentPage, onNavClick}) => {

  var height = {height: window.innerHeight-100}


          // <Route path="/">
          // </Route>
            // <Route path="/" component={MapPage} />
            // <Route path="/journal" component={JournalPage}/>
            // <Route path="/data" component={DataPage}/>
            // <Route path="/about" component={AboutPage}/>
            // <Route path="/share" component={SharePage}/>
            // <MapPage active={true}/>
            // <JournalPage active={false}/>
            // <DataPage active={false}/>
            // <AboutPage active={false}/>
            // <SharePage active={false}/>
  return (
    <div>
      <Navigation currentPage={currentPage}/>
      <div id="content" style={height}>
        <LightBox active={false}/>
        <Timeline/>
        <div id="pageContainer">
        </div>
      </div>
    </div>
  )
}

Content.propTypes = {
  currentPage : PropTypes.string.isRequired
}

export default Content
