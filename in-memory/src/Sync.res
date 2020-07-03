open Model

type handler = t => (t, effects)

let incrementTicked: handler =
  current => {
    ({...current, ticked: current.ticked + 1}, list[])
  }

// TODO: fetches again the changes feed after some time
let checkFetchChanges: handler =
  model => {
    switch model.states.changes {
    | ChangesFeedNeverFetched => (
        {
          ...model,
          states: {
            changes: ChangesFeedCurrentlyFetching,
          },
        },
        list[FetchChangesFeed],
      )
    | _ => (model, list[])
    }
  }

let combineHandlers = (handlers: list<handler>): handler => {
  let rec aux = (handlers: list<handler>, model: t, acc: effects) => {
    switch handlers {
    | list[] => (model, acc)
    | list[handler, ...rest] =>
      let (next, effects) = handler(model)
      aux(rest, next, List.concat(acc, effects))
    }
  }
  model => aux(handlers, model, list[])
}

let handleTick = combineHandlers(list[incrementTicked, checkFetchChanges])

let update = (model: t, event: event): (t, effects) => {
  switch event {
  | Tick => handleTick(model)
  | _ => (model, list[])
  }
}
