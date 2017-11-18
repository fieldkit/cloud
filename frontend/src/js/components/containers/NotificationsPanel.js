// @flow weak

import React, { Component } from 'react';

import ReactCSSTransitionGroup from 'react-addons-css-transition-group';

import FeaturePanel from './FeaturePanel';

const containerStyle: React.CSSProperties = {
    position: 'absolute',
    width: '300px',
    zIndex: 10,
    display: 'flex',
    flexDirection: 'column',
};

const panelStyle: React.CSSProperties = {
    boxShadow: '0px 1px 4px rgba(0, 0, 0, .3)',
    border: '1px solid rgba(0, 0, 0, 0.4)'
};

const separatorStyle: React.CSSProperties = {
    marginTop: '10px'
};

type Props = {
    geojson: GeoJSON
};

export default class NotificationsPanel extends Component {
    props: Props

    render() {
        const { features } = this.props;
        const position = { top: 150, right: 30, bottom: 'auto', left: 'auto' };

        if (features.length === 0) {
            return (<div></div>);
        }

        const panels = features.map((f, i) => <FeaturePanel key={i} feature={f} style={{ ...panelStyle, ...(i > 0 ? separatorStyle : {}) }} />);

        return (
            <div className="notification-panel" style={{ ...containerStyle, ...position }}>
                <ReactCSSTransitionGroup transitionName="transition" transitionEnter={true} transitionEnterTimeout={500} transitionLeave={true}
                    transitionLeaveTimeout={500}>
                    {panels}
                </ReactCSSTransitionGroup>
            </div>
        )
    }
}
