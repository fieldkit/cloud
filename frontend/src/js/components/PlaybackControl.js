import React, { Component } from 'react';

const containerStyle: React.CSSProperties = {
    position: 'absolute',
    zIndex: 10,
    display: 'flex',
    flexDirection: 'column',
    boxShadow: '0px 1px 4px rgba(0, 0, 0, .3)',
    border: '1px solid rgba(0, 0, 0, 0.1)'
};

const buttonContainerStyle: React.CSSProperties = {
    display: 'flex',
    flexDirection: 'row',
};

const advancedContainerStyle: React.CSSProperties = {
    backgroundColor: '#f9f9f9',
    padding: "5px",
    color: "#000",
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

const buttonActiveStyle = {
    backgroundColor: '#d0d0d0',
};

const buttonLastStyle = {
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
    state: {
        advanced: boolean
    }

    constructor() {
        super();
        this.state = {
            advanced: false
        };
    }

    onToggleView() {
        this.setState({
            advanced: !this.state.advanced
        });
    }

    renderAdvanced() {
        return (
            <div style={{ ...advancedContainerStyle }}>
                Advanced Options
            </div>
        );
    }

    render() {
        const { className, playback, onPlaybackChange } = this.props;
        const { advanced } = this.state;

        const position = { top: 100, right: 30, bottom: 'auto', left: 'auto' };
        const smallStyle = {
            ...buttonStyle,
            ...{ width: '30px' }
        };

        const playOrPause = (playback === PlaybackModes.Pause ||
                             playback === PlaybackModes.Beginning ||
                             playback === PlaybackModes.End) ? PlaybackModes.Play : PlaybackModes.Pause;
        const playOrPauseIcon = playOrPause === PlaybackModes.Play ? "fa fa-play" : "fa fa-pause";
        // <i class="fa fa-fast-backward" aria-hidden="true"></i>
        // <i class="fa fa-fast-forward" aria-hidden="true"></i>

        return (
            <div className={className} style={{ ...containerStyle, ...position }}>
                <div style={{ ...buttonContainerStyle }}>

                    <button type="button" style={{ ...smallStyle, ...buttonInnerStyle, ...buttonTopStyle  }} onClick={ onPlaybackChange.bind(this, PlaybackModes.Beginning) }><i className="fa fa-backward" aria-hidden="true"></i></button>
                    <button type="button" style={{ ...smallStyle, ...buttonInnerStyle }} onClick={ onPlaybackChange.bind(this, PlaybackModes.Rewind) }><i className="fa fa-step-backward" aria-hidden="true"></i></button>
                    <button type="button" style={{ ...smallStyle, ...buttonInnerStyle, ...buttonActiveStyle }} onClick={ onPlaybackChange.bind(this, playOrPause) }><i className={playOrPauseIcon} aria-hidden="true"></i></button>
                    <button type="button" style={{ ...smallStyle, ...buttonInnerStyle }} onClick={ onPlaybackChange.bind(this, PlaybackModes.Forward) }><i className="fa fa-forward" aria-hidden="true"></i></button>
                    <button type="button" style={{ ...smallStyle, ...buttonInnerStyle }} onClick={ onPlaybackChange.bind(this, PlaybackModes.End) }><i className="fa fa-step-forward" aria-hidden="true"></i></button>
                    <button type="button" style={{ ...smallStyle, ...buttonLastStyle }} onClick={ this.onToggleView.bind(this) }><i className="fa fa-ellipsis-h" aria-hidden="true"></i></button>
                </div>
                { advanced && this.renderAdvanced() }
            </div>
        )
    }
}
