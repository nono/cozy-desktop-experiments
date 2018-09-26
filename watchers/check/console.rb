#!/usr/bin/env ruby

require_relative 'lib/boot'

at_exit { Helpers.cleanup }
Pry.start binding, prompt: Pry::SIMPLE_PROMPT, quiet: true
