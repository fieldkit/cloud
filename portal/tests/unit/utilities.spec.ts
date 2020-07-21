import * as utils from "@/utilities";

describe("Utilities", () => {
    it("Converts bytes to a label", () => {
        expect(utils.convertBytesToLabel(453453453)).toBe("432 MB");
    });
});
