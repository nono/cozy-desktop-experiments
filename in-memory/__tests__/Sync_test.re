open Jest;

describe("Sync", () => {
  Expect.(
    describe("update", () => {
      test("tick", () => {
        let (next, _) =
          Model.init({cozyURL: "http://cozy.tools:8080/"})
          ->Sync.update(Model.Tick);
        expect(next.ticked) |> toBe(1);
      })
    })
  )
});
