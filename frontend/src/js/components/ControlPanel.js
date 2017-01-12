
import React, {PropTypes} from 'react'
import YearSelector from './YearSelector'
import DateSelector from './DateSelector'
import PlaybackSelector from './PlaybackSelector'
import FocusSelector from './FocusSelector'
import ZoomSelector from './ZoomSelector'
import LayoutSelector from './LayoutSelector'

const ControlPanel = ({pathName, expeditionID, expeditions, currentDate, playback, mainFocus, secondaryFocus, zoom, layout, onYearChange, onPlaybackChange, onMainFocusChange, onSecondaryFocusChange, onZoomChange, onLayoutChange, viewport}) => {
  if (!expeditionID) return <div className="controlPanel"></div>

  // {pathName === '/map' ? <FocusSelector mainFocus={mainFocus} secondaryFocus={secondaryFocus} onMainFocusChange={onMainFocusChange} onSecondaryFocusChange={onSecondaryFocusChange}/> : null}
  // {pathName === '/journal' ? <LayoutSelector mode={layout} onLayoutChange={onLayoutChange}/> : null}
      // <YearSelector expeditions={expeditions} expeditionID={expeditionID} onYearChange={onYearChange}/>

  return (
    <div className="controlPanel">
      <DateSelector expeditions={expeditions} expeditionID={expeditionID} currentDate={currentDate} />
      {location.pathname.indexOf('/map') > -1 && <PlaybackSelector mode={playback} onPlaybackChange={onPlaybackChange}/>}
      {location.pathname.indexOf('/map') > -1 && window.innerWidth > 768 && <ZoomSelector onZoomChange={onZoomChange} viewport={viewport}/>}
    </div>
  )
}

ControlPanel.propTypes = {
  expeditionID: PropTypes.string,
  expeditions: PropTypes.object,
  currentDate: PropTypes.object,
  playback: PropTypes.string,
  mainFocus: PropTypes.string,
  secondaryFocus: PropTypes.string,
  zoom: PropTypes.number,
  layout: PropTypes.string,
  onYearChange: PropTypes.func.isRequired,
  onPlaybackChange: PropTypes.func.isRequired,
  onMainFocusChange: PropTypes.func.isRequired,
  onSecondaryFocusChange: PropTypes.func.isRequired,
  onZoomChange: PropTypes.func.isRequired,
  onLayoutChange: PropTypes.func.isRequired,
  viewPort: PropTypes.object
}

export default ControlPanel
