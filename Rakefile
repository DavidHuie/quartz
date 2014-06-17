require 'rubygems'
require 'rake'
require 'jeweler'

Jeweler::Tasks.new do |gem|
  gem.name = 'quartz'
  gem.homepage = 'http://github.com/DavidHuie/quartz'
  gem.license = 'MIT'
  gem.summary = 'A gem for calling Go code from Ruby'
  gem.description = 'A gem for calling Go code from Ruby'
  gem.email = 'dahuie@gmail.com'
  gem.authors = ['David Huie']
end

Jeweler::RubygemsDotOrgTasks.new

require 'rspec/core'
require 'rspec/core/rake_task'

RSpec::Core::RakeTask.new(:spec) do |spec|
  spec.pattern = FileList['spec/**/*_spec.rb']
end

desc 'Code coverage detail'
task :simplecov do
  ENV['COVERAGE'] = '1'
  Rake::Task['spec'].execute
end

task :default => :spec
