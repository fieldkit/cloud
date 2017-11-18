import React, { Component } from 'react';

const containerStyle: React.CSSProperties = {
    position: 'absolute',
    zIndex: 10,
    display: 'flex',
    flexDirection: 'column',
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
    height: 26,
    width: 26,
    backgroundPosition: '0px 0px',
    backgroundSize: '26px 260px',
    outline: 0
};

const buttonTopStyle = {
    borderTopLeftRadius: 2,
    borderTopRightRadius: 2
};

const buttonInnerStyle = {
    borderBottom: '1px solid rgba(0, 0, 0, 0.1)',
};

const buttonBottomStyle = {
    borderBottomLeftRadius: 2,
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

        const position = { top: 200, left: 30, bottom: 'auto', right: 'auto' };
        const smallStyle = {
            ...buttonStyle,
            ...{ width: '80px' }
        };

        const playOrPause = (playback === PlaybackModes.Pause ||
                             playback === PlaybackModes.Beginning ||
                             playback === PlaybackModes.End) ? PlaybackModes.Play : PlaybackModes.Pause;

        return (
            <div
                className={className}
                style={{ ...containerStyle, ...position }}>

                <button type="button" style={{ ...smallStyle, ...buttonInnerStyle, ...buttonTopStyle  }} onClick={ onPlaybackChange.bind(this, PlaybackModes.Beginning) }>Beg</button>
                <button type="button" style={{ ...smallStyle, ...buttonInnerStyle }} onClick={ onPlaybackChange.bind(this, PlaybackModes.Rewind) }>RW</button>
                <button type="button" style={{ ...smallStyle, ...buttonInnerStyle }} onClick={ onPlaybackChange.bind(this, playOrPause) }>{ playOrPause.mode }</button>
                <button type="button" style={{ ...smallStyle, ...buttonInnerStyle }} onClick={ onPlaybackChange.bind(this, PlaybackModes.Forward) }>FF</button>
                <button type="button" style={{ ...smallStyle, ...buttonBottomStyle }} onClick={ onPlaybackChange.bind(this, PlaybackModes.End) }>End</button>
            </div>
        )
    }
}
