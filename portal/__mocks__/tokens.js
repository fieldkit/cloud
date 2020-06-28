const mock = jest.fn().mockImplementation(() => {
    return {
        getToken: () => null,
    };
});

export default mock;
