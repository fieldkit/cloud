import { expect, config } from 'chai';
import { avgSelection, countSelection, simpleStringSelection, simpleNumSelection, makeSelection, equalGrouping, getGroupingFn, transform, Broker, wholeGroupSelection, simpleLocationSelection, generateSelectionsFromMap } from '../lib/caboodle' 
import {Location} from '../lib/proto/flow'

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
  
  it('Simple Strings Correctly', () => {
    const test_selection = {
      id: 1,
      value_name: "test_value",
      source_attribute: "source",
      operation: "simple_string"
    }
    const test_data = [{source: "foo"}, {source: "bar"}, {source: "baz"}]
    const test_data_2 = [{source: 1}, {source: "bar"}, {source: "baz"}]
    const test_data_3 = []
    const selection_fn = simpleStringSelection(test_selection)

    expect(selection_fn(test_data)).to.equal("foo");
    expect(selection_fn(test_data_2)).to.equal("bar");
    expect(selection_fn(test_data_3)).to.equal("");
  });
  
  it('Simple Nums Correctly', () => {
    const test_selection = {
      id: 1,
      value_name: "test_value",
      source_attribute: "source",
      operation: "simple_num"
    }
    const test_data = [{source: 1}, {source: 2}, {source: 3}]
    const test_data_2 = [{source: "1"}, {source: 2}, {source: 3}]
    const test_data_3 = []
    const selection_fn = simpleNumSelection(test_selection)

    expect(selection_fn(test_data)).to.equal(1);
    expect(selection_fn(test_data_2)).to.equal(2);
    expect(selection_fn(test_data_3)).to.equal(0);
  });
  
  it('Simple Locs Correctly', () => {
    const test_selection = {
      id: 1,
      value_name: "test_value",
      source_attribute: "source",
      operation: "simple_location"
    }
    const loc = new Location({latitude: 3, longitude: 4})
    const empty_loc = new Location({latitude: 0, longitude: 0, altitude: 0})
    const test_data = [{source: loc}, {source: "bar"}, {source: "baz"}]
    const test_data_2 = [{source: 1}, {source: loc}, {source: "baz"}]
    const test_data_3 = []
    const selection_fn = simpleLocationSelection(test_selection)

    expect(selection_fn(test_data)).to.deep.equal(loc);
    expect(selection_fn(test_data_2)).to.deep.equal(loc);
    expect(selection_fn(test_data_3)).to.deep.equal(empty_loc);
  });

  it('Performs a Whole Group Selection', () => {
    const test_selection = {
      id: 1,
      value_name: "test_value",
      source_attribute: "",
      operation: "whole_group"
    }
    const test_data = [{source: 1}, {source: 2}, {source: 3}]
    const selection_fn = wholeGroupSelection(test_selection)

    expect(selection_fn(test_data)).to.deep.equal(test_data);
  })
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

  it('Correctly Generates a Simple Selection', () => {
    const test_selection = generateSelectionsFromMap({
      "foo": "number",
      "bar": "string",
      "baz": "location"
    })

    const expectedResult = [
      {id: 0, value_name: "foo", source_attribute: "foo", operation: "simple_num"},
      {id: 1, value_name: "bar", source_attribute: "bar", operation: "simple_string"},
      {id: 2, value_name: "baz", source_attribute: "baz", operation: "simple_location"}
    ]

    expect(expectedResult).to.deep.equal(test_selection)
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
