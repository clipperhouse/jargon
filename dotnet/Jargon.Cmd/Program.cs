﻿using CommandLine;
using System;
using System.Diagnostics;
using System.IO;

namespace Jargon.Cmd
{
    public sealed class Program
    {
        private sealed class Options
        {
            [Option('f', Required = false, HelpText = "Input file path")]
            public string InputFile { get; set; }

            internal bool FileIsSet => !string.IsNullOrEmpty(InputFile);

            [Option('s', Required = false, HelpText = "A (quoted) string to lemmatize")]
            public string String { get; set; }

            internal bool StringIsSet => !string.IsNullOrEmpty(String);

            [Option('u', Required = false, HelpText = "A URL to fetch and lemmatize")]
            public string Url { get; set; }

            internal bool UrlIsSet => !string.IsNullOrEmpty(Url);

            [Option('o', Required = false, HelpText = "Output file path")]
            public string OutputFile { get; set; }

            internal bool OutputFileIsSet => !string.IsNullOrEmpty(OutputFile);
        }

        public static void Main(string[] args)
        {
            var res = Parser.Default.ParseArguments<Options>(args);
            res.WithParsed(HandleOptions).WithNotParsed(e => PrintExamplesAndExit());
        }

        private static void HandleOptions(Options opts)
        {
            var numInputsSet =
                (opts.FileIsSet ? 1 : 0) +
                (opts.StringIsSet ? 1 : 0) +
                (opts.UrlIsSet ? 1 : 0);
            if (numInputsSet > 1)
            {
                Console.WriteLine($"Only one of `f`, `s`, and `u` may be set");
                Environment.Exit(-2);
            }

            TextWriter writer = null;
            try
            {
                if (opts.OutputFileIsSet)
                {
                    try
                    {
                        var fs = File.Create(opts.OutputFile);
                        writer = new StreamWriter(fs);
                    }
                    catch (Exception e)
                    {
                        Console.WriteLine($"Could not create file ({e.Message}): {opts.OutputFile}");
                        Environment.Exit(-5);
                    }
                }
                else
                {
                    writer = Console.Out;
                }

                if (opts.FileIsSet)
                {
                    if (!File.Exists(opts.InputFile))
                    {
                        Console.WriteLine($"Could not find file: {opts.InputFile}");
                        Environment.Exit(-3);
                    }

                    try
                    {
                        using (var reader = new StreamReader(File.OpenRead(opts.InputFile)))
                        {
                            Lemmatize(reader, writer);
                            return;
                        }
                    }
                    catch (Exception e)
                    {
                        Console.WriteLine($"Could not read file ({e.Message}): {opts.InputFile}");
                        Environment.Exit(-4);
                    }
                }

                if (opts.StringIsSet)
                {
                    using (var reader = new StringReader(opts.String))
                    {
                        Lemmatize(reader, writer);
                        return;
                    }
                }

                if (opts.UrlIsSet)
                {
                    try
                    {
                        using (var web = new CompressedWebClient())
                        {
                            var html = web.DownloadString(opts.Url);
                            using (var reader = new StringReader(html))
                            {
                                Lemmatize(reader, writer);
                            }
                            return;
                        }
                    }
                    catch (Exception e)
                    {
                        Console.WriteLine($"Could not download url ({e.Message}): {opts.Url}");
                        Environment.Exit(-5);
                    }
                }

                Lemmatize(Console.In, writer);
            }
            finally
            {
                if (writer != null && !object.ReferenceEquals(writer, Console.Out))
                {
                    writer.Dispose();
                }
            }
        }

        private static void Lemmatize(TextReader reader, TextWriter writer)
        {
            var lemmatizer = new Lemmatizer(Data.StackExchange.Instance, 3);
            using(var toks = new TextTokens(reader))
            using(var e = new LemmaTokens(lemmatizer, toks))
            {
                while (e.MoveNext())
                {
                    writer.Write(e.Current.Value);
                }
            }
        }

        private static void PrintExamplesAndExit()
        {
            string cmdName;
            using (var proc = Process.GetCurrentProcess())
            {
                cmdName = proc.ProcessName;
            }

            var pathSeparator = Path.DirectorySeparatorChar;

            Console.WriteLine($@"
Usage: {cmdName} accepts piped UTF8 text from tools such as cat, curl or echo, via Stdin
		
  Example: echo ""I luv Rails"" | {cmdName}
Alternatively, use {cmdName} 'standalone' by passing flags for inputs and outputs:
  Example: {cmdName} -f {pathSeparator}path{pathSeparator}to{pathSeparator}original.txt -o {pathSeparator}path{pathSeparator}to{pathSeparator}lemmatized.txt

Results are piped to Stdout (regardless of input)");

            Environment.Exit(-1);
        }
    }
}
