import station from "../tests/unit/fixtures/station";
const user = { bio: "A Test user.", email: "test@conservify.org", id: 1, name: "Test User" };

const mock = jest.fn().mockImplementation(() => {
    return {
        getCurrentUser: () => Promise.resolve(user),
        getStation: () => Promise.resolve(station),
        getStations: () => Promise.resolve({ stations: [station] }),
        getUserProjects: () => Promise.resolve({ projects: [] }),
        getModulesMeta: () => Promise.resolve([]),
    };
});

export default mock;
