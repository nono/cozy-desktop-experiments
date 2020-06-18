open Jest;

describe("model", () => {
  Expect.(
    describe("update", () => {
      test("tick", () => {
        let (next, _) =
          Model.init({cozyURL: "http://cozy.tools:8080/"})
          ->Model.update(Event.Tick);
        expect(next.ticked) |> toBe(1);
      })
    })
  )
});
