type property
type generator<'a>
type predicate<'a> = 'a => bool

@bs.module("fast-check") external string: unit => generator<string> = "string"
@bs.module("fast-check") external array: generator<'a> => generator<array<'a>> = "array"

@bs.module("fast-check") external property: (generator<'a>, predicate<'a>) => property = "property"
@bs.module("fast-check") external _assert: property => unit = "assert"

let expect = (generator: generator<'a>, predicate: predicate<'a>) => {
  Jest.Expect.expect(() => {
    _assert(property(generator, predicate))
  }) |> Jest.Expect.not_ |> Jest.Expect.toThrow
}
