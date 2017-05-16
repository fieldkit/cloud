import { expect } from 'chai';
import {foo} from '../lib/caboodle' 

describe('Arithmetic', () => {  
  it('should calculate 1 + 1 correctly', () => {
    const expectedResult = 2;

    expect(foo(1,1)).to.equal(expectedResult);
  });
});
