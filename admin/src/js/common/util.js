// @flow
import React from 'react';
import slug from 'slug';

export class BaseError {
  name: string;
  message: string;
  stack: ?string;

  constructor(message: string = 'Error') {
    this.name = this.constructor.name;
    this.message = message;
    if (typeof Error.captureStackTrace === 'function') {
      Error.captureStackTrace(this, this.constructor);
    } else {
      this.stack = (new Error(message)).stack;
    }
  }
}

export type ErrorMap = { [key: string]: string[] };

export function errorsFor(errors: ?ErrorMap, key: string) {
  if (errors && errors[key]) {
    return errors[key].map((error, i) =>  <p key={`errors-${key}-${i}`} className="error">{error}</p>);
  }
}

export function slugify(input: string): string {
  return slug(input).toLowerCase();
}
