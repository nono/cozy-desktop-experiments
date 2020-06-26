open Jest;

describe("Sync", () => {
  Expect.(
    describe("update", () => {
      test("tick", () => {
        let (updated, _) =
          Model.init({cozyURL: "http://cozy.tools:8080/"})
          ->Sync.update(Model.Tick);
        expect(updated.ticked) |> toBe(1);
      })
    })
  )
});
