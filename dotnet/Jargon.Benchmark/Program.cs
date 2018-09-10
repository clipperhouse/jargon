using BenchmarkDotNet.Configs;
using BenchmarkDotNet.Diagnosers;
using BenchmarkDotNet.Jobs;
using BenchmarkDotNet.Running;
using System.Linq;
using System.Reflection;

namespace Jargon.Benchmark
{
    class Program
    {
        static void Main(string[] args)
        {
            //var test = new Benchmarks.TokenizeBench();
            //test.LoadData();

            //for (var i = 0; i < 10000; i++)
            //{
            //    test.TokenizeBenchmark();
            //}

            var config = ManualConfig.CreateEmpty().With(new MemoryDiagnoser()).With(DefaultConfig.Instance.GetColumnProviders().ToArray()).With(DefaultConfig.Instance.GetExporters().ToArray());
            config = config.With(Job.RyuJitX64);
            config = config.With(DefaultConfig.Instance.GetLoggers().ToArray());

            BenchmarkRunner.Run(Assembly.GetExecutingAssembly(), config);
        }
    }
}
