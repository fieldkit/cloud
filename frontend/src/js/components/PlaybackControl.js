import React, { Component } from 'react';

const containerStyle: React.CSSProperties = {
    position: 'absolute',
    zIndex: 10,
    display: 'flex',
    flexDirection: 'row',
    boxShadow: '0px 1px 4px rgba(0, 0, 0, .3)',
    border: '1px solid rgba(0, 0, 0, 0.1)'
};

const buttonStyle = {
    backgroundColor: '#f9f9f9',
    opacity: 0.95,
    transition: 'background-color 0.16s ease-out',
    cursor: 'pointer',
    fontWeight: "bold",
    border: 0,
    height: 32,
    width: 36,
    backgroundPosition: '0px 0px',
    backgroundSize: '26px 260px',
    outline: 0
};

const buttonTopStyle = {
    borderTopLeftRadius: 2,
    borderBottomLeftRadius: 2
};

const buttonInnerStyle = {
    borderRight: '1px solid rgba(0, 0, 0, 0.1)',
};

const buttonBottomStyle = {
    borderTopRightRadius: 2,
    borderBottomRightRadius: 2
};


export const PlaybackModes = {
    Beginning: { mode: 'Beginning' },
    Rewind: { mode: 'Rewind' },
    Play: { mode: 'Play' },
    Pause: { mode: 'Pause' },
    Forward: { mode: 'Forward' },
    End: { mode: 'End' },
};

export default class PlaybackControl extends Component {
    props: {
        className: PropTypes.string,
        playback: PropTypes.object.isRequired,
        onPlaybackChange: PropTypes.func.isRequired
    }

    render() {
        const { className, playback, onPlaybackChange } = this.props;

        const position = { top: 100, right: 30, bottom: 'auto', left: 'auto' };
        const smallStyle = {
            ...buttonStyle,
            ...{ width: '30px' }
        };

        const buttonClasses = "button";
        const playOrPause = (playback === PlaybackModes.Pause ||
                             playback === PlaybackModes.Beginning ||
                             playback === PlaybackModes.End) ? PlaybackModes.Play : PlaybackModes.Pause;
        const playOrPauseIcon = playOrPause == PlaybackModes.Play ? "fa fa-play" : "fa fa-pause";

        return (
            <div
                className={className}
                style={{ ...containerStyle, ...position }}>

                <button type="button" style={{ ...smallStyle, ...buttonInnerStyle, ...buttonTopStyle  }} onClick={ onPlaybackChange.bind(this, PlaybackModes.Beginning) }><i className="fa fa-backward" aria-hidden="true"></i></button>
                <button type="button" style={{ ...smallStyle, ...buttonInnerStyle }} onClick={ onPlaybackChange.bind(this, PlaybackModes.Rewind) }><i className="fa fa-step-backward" aria-hidden="true"></i></button>
                <button type="button" style={{ ...smallStyle, ...buttonInnerStyle }} onClick={ onPlaybackChange.bind(this, playOrPause) }><i className={playOrPauseIcon} aria-hidden="true"></i></button>
                <button type="button" style={{ ...smallStyle, ...buttonInnerStyle }} onClick={ onPlaybackChange.bind(this, PlaybackModes.Forward) }><i className="fa fa-forward" aria-hidden="true"></i></button>
                <button type="button" style={{ ...smallStyle, ...buttonBottomStyle }} onClick={ onPlaybackChange.bind(this, PlaybackModes.End) }><i className="fa fa-step-forward" aria-hidden="true"></i></button>
            </div>
        )
    }
}
