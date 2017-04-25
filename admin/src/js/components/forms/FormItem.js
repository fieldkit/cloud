// @flow

import React, { Component } from 'react';

import { errorsFor } from '../../common/util';
import type { APIErrors } from '../../api/types';

type Props = {
  labelText: string;
  name: string;
  inline?: boolean;
  className?: string;
  type?: string;
  ref?: string;
  value: string | number;
  errors: ?APIErrors;
  onChange?: (e: any) => void;
}

export class FormItem extends Component<void, $Exact<Props>, void> {
  render() {
    const {
      name,
      inline,
      className,
      labelText,
      type,
      ref,
      value,
      errors,
      onChange
    } = this.props;

    return (
      <div className={inline ? 'form-group inline' : 'form-group' }>
        <label htmlFor={name}>{labelText}</label>
        <input ref={ref} name={name} type={type || 'text'} value={value} onChange={onChange}/>
        { errorsFor(errors, name) }
      </div>
    );
  }
}
