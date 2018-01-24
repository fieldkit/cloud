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

    positionStyle() {
        const { sidePanelVisible } = this.props;
        if (!sidePanelVisible) {
            return { top: 150, right: 30, bottom: 'auto', left: 'auto' };
        }
        return { top: 150, right: 30 + 300, bottom: 'auto', left: 'auto' };
    }


    render() {
        const { features } = this.props;
        const position = this.positionStyle();

        if (features.length === 0) {
            return (<div></div>);
        }

        const visible = features.slice(0, 5);

        const panels = visible.map((f, i) => <FeaturePanel key={i} feature={f} style={{ ...panelStyle, ...(i > 0 ? separatorStyle : {}) }} />);

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
