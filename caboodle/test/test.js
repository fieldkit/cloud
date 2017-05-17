import { expect, config } from 'chai';
import { avgSelection, countSelection, makeSelection, equalGrouping, getGroupingFn, transform, Broker } from '../lib/caboodle' 

describe('Selection Functions', () => {  
  it('Average Groups Correctly', () => {
    const expectedResult = 2;
    const test_selection = {
      id: 1,
      value_name: "test_value",
      source_attribute: "source",
      operation: "avg"
    }
    const test_data = [{source: 1}, {source: 2}, {source: 3}]
    const selection_fn = avgSelection(test_selection)

    expect(selection_fn(test_data)).to.equal(expectedResult);
  });
  
  it('Counts Groups Correctly', () => {
    const expectedResult = 3;
    const test_selection = {
      id: 1,
      value_name: "test_value",
      source_attribute: "source",
      operation: "count"
    }
    const test_data = [{source: 1}, {source: 2}, {source: 3}]
    const selection_fn = countSelection(test_selection)

    expect(selection_fn(test_data)).to.equal(expectedResult);
  });
});

describe('Selections', () => {
  it('Correctly Makes Selections',() => {
    const test_selection = {
      id: 1,
      value_name: "test_value",
      source_attribute: "source",
      operation: "avg"
    }
    const selection_tx = makeSelection(test_selection) 
    let o_data = [{source: 1}, {source: 2}, {source: 3}]
    let {data, output} = selection_tx({data: o_data, output: {}})
    const expectedResult = {test_value: 2}
    expect(output).to.deep.equal(expectedResult)
  })
})

describe('Grouping Functions', () => {  
  it('Groups Equal Correctly', () => {
    const test_grouping = {
      operation: "equal",
      parameter: null,
      source_attribute: "id"
    }
    const grouping_fn = equalGrouping(test_grouping)
    const test_data = [
      [{id: 1, value: 0}, {id: 2, value: 2},{id: 3, value: 4}],
      [{id: 1, value: 1}, {id: 4, value: 6},{id: 2, value: 7}]
    ]
    const groups = grouping_fn(test_data)
    const expectedResult = [
      [{id: 1, value: 0},{id: 1, value: 1}],
      [{id: 2, value: 2},{id: 2, value: 7}],
      [{id: 3, value: 4}],
      [{id: 4, value: 6}]
    ]
    expect(groups).to.deep.equal(expectedResult);
  });
});

describe('Grouping', () => {
  it('Correctly Makes Groups', () => {
    const test_grouping = {
      operation: "equal",
      parameter: null,
      source_attribute: "id"
    }
    const grouping_fn = getGroupingFn(test_grouping)
    const test_data = [
      [{id: 1, value: 0}, {id: 2, value: 2},{id: 3, value: 4}],
      [{id: 1, value: 1}, {id: 4, value: 6},{id: 2, value: 7}]
    ]
    const groups = grouping_fn(test_data)
    const expectedResult = [
      [{id: 1, value: 0},{id: 1, value: 1}],
      [{id: 2, value: 2},{id: 2, value: 7}],
      [{id: 3, value: 4}],
      [{id: 4, value: 6}]
    ]
    expect(groups).to.deep.equal(expectedResult);
  })
})

describe('Transformation', () => {
  it('Correctly Transforms Data', () => {
    const test_grouping = {
      operation: "equal",
      parameter: null,
      source_attribute: "id"
    }
    const test_selection = {
      id: 1,
      value_name: "test_value",
      source_attribute: "value",
      operation: "avg"
    }
    const test_data = [
      [{id: 1, value: 0}, {id: 2, value: 2},{id: 3, value: 4}],
      [{id: 1, value: 1}, {id: 4, value: 6},{id: 2, value: 7}]
    ]
    const final_stream = transform(test_data, test_grouping, [test_selection])
    const expectedResult = [{test_value: 0.5},{test_value: 4.5},{test_value:4},{test_value:6}]

    expect(final_stream).to.deep.equal(expectedResult);
  })
})

describe('Broker', () => {
  it('Prevents Restarting a Broker', () => {
    const broker = new Broker("")       
    broker.start()
    expect(broker.start).to.throw(Error)
  })

  it('Prevent Registering After a Start', () => {
    const test_grouping = {
      operation: "equal",
      parameter: null,
      source_attribute: "id"
    }
    const test_selection = {
      id: 1,
      value_name: "test_value",
      source_attribute: "value",
      operation: "avg"
    }
    const broker = new Broker("")

    expect(() => {
      broker.start()
    }).not.to.throw(Error)

    expect(() => {
      broker.register([],test_grouping,[test_selection],()=>{})
    }).to.throw(Error)
  })
  
  it('Handles a Response Correctly',() => {
    const test_grouping = {
      operation: "equal",
      parameter: null,
      source_attribute: "id"
    }
    const test_selection = {
      id: 1,
      value_name: "test_value",
      source_attribute: "value",
      operation: "avg"
    }
    const test_data = [
      [{id: 1, value: 0}, {id: 2, value: 2},{id: 3, value: 4}],
      [{id: 1, value: 1}, {id: 4, value: 6},{id: 2, value: 7}]
    ]
    const broker = new Broker("")
    const registration = broker.register([],test_grouping,[test_selection],()=>{})
    
    const result = broker.handleResponse(registration.id, {streams: test_data},registration.processor,() => {})

    const expectedResult = [{test_value: 0.5},{test_value: 4.5},{test_value:4},{test_value:6}]

    expect(result).to.deep.equal(expectedResult);
  })
})
