#!/usr/bin/env node

const Ignore = require('desktop/core/ignore')
const Merge = require('desktop/core/merge')
const Pouch = require('desktop/core/pouch')
const Prep = require('desktop/core/prep')
const Watcher = require('desktop/core/local/watcher')

const syncPath = '../workspace' // TODO
const config = {
  dbPath: '../tmp' // TODO
}
const pouch = new Pouch(config)

const ignore = new Ignore('')
const merge = new Merge(pouch)
const prep = new Prep(merge, ignore, config)

const loggerPrep = {
  addFileAsync: (side, doc) => {
    console.log('File +', doc)
    return prep.addFileAsync(side, doc)
  },
  moveFileAsync: (side, doc, old) => {
    console.log('File <', old, '>', doc)
    return prep.moveFileAsync(side, doc, old)
  },
  moveFolderAsync: (side, doc, old) => {
    console.log('Dir  <', old, '>', doc)
    return prep.moveFolderAsync(side, doc, old)
  },
  putFolderAsync: (side, doc) => {
    console.log('Dir  +', doc)
    return prep.putFolderAsync(side, doc)
  },
  trashFileAsync: (side, doc) => {
    console.log('File -', doc)
    return prep.trashFileAsync(side, doc)
  },
  trashFolderAsync: (side, doc) => {
    console.log('Dir  -', doc)
    return prep.trashFolderAsync(side, doc)
  },
  updateFileAsync: (side, doc) => {
    console.log('File *', doc)
    return prep.updateFileAsync(side, doc)
  },
}

const events = {
  emit: (msg) => {
    console.log('Info', msg)
  }
}

pouch.addAllViews(() => {
  const w = new Watcher(syncPath, loggerPrep, pouch, events)
  w.start()
  const quit = () => {
    w.stop()
  }
  process.on('SIGINT', quit)
  process.on('SIGTERM', quit)
})
