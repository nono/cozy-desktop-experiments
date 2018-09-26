#!/usr/bin/env node

const Pouch = require('desktop/core/pouch')
const Watcher = require('desktop/core/local/watcher')

const syncPath = '../workspace' // TODO

const config = {
  dbPath: '../tmp' // TODO
}

const prep = {
  addFileAsync: (side, doc) => {
    console.log('F +', doc)
    return Promise.resolve()
  },
  moveFileAsync: (side, doc, old) => {
    console.log('F <', old, '>', doc)
    return Promise.resolve()
  },
  moveFolderAsync: (side, doc, old) => {
    console.log('D <', old, '>', doc)
    return Promise.resolve()
  },
  putFolderAsync: (side, doc) => {
    console.log('D +', doc)
    return Promise.resolve()
  },
  trashFileAsync: (side, doc) => {
    console.log('F -', doc)
    return Promise.resolve()
  },
  trashFolderAsync: (side, doc) => {
    console.log('D -', doc)
    return Promise.resolve()
  },
  updateFileAsync: (side, doc) => {
    console.log('F *', doc)
    return Promise.resolve()
  },
}

const events = {
  emit: (msg) => {
    console.log('E', msg)
  }
}

const pouch = new Pouch(config)
pouch.addAllViews(() => {
  const w = new Watcher(syncPath, prep, pouch, events)
  w.start()
  const quit = () => {
    w.stop()
  }
  process.on('SIGINT', quit)
  process.on('SIGTERM', quit)
})
