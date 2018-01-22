// @flow weak

import React, { Component } from 'react';

const mainStyle: React.CSSProperties = {
    backgroundColor: '#f9f9f9',
    color: "#000",
    padding: "10px",
    fontWeight: "bold",
    fontSize: "14px",
};

export class Loading extends Component {
    render() {
        return (
            <div style={mainStyle}>
                Loading
            </div>
        );
    }
}

