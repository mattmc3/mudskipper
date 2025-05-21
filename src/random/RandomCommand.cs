using System;
using System.Collections.Generic;
using System.Linq;

namespace Mudskipper;

public class RandomCommand
{
    public static void Main(string[] args)
    {
        if (args.Contains("--help") || args.Contains("-h"))
        {
            Usage();
            return;
        }

        try
        {
            Run(args);
        }
        catch (Exception ex)
        {
            Console.Error.WriteLine($"random: error: {ex.Message}");
            Environment.Exit(1);
        }
    }

    public static void Run(string[] args)
    {
        if (args.Length == 0)
        {
            // Default behavior: Print a random number between 0 and 32767
            Console.WriteLine(Random.Shared.Next(0, 32768));
            return;
        }

        if (args[0] == "choice")
        {
            // random choice subcommand: Output a random selection from the specified arguments
            if (args.Length == 1)
            {
                Console.Error.WriteLine("random: nothing to choose from");
                Environment.Exit(1);
            }

            string selection = args[Random.Shared.Next(1, args.Length)];
            Console.WriteLine(selection);
            return;
        }

        // Make sure everything is an integer, or tell the user
        foreach (var arg in args)
        {
            if (!int.TryParse(arg, out _))
            {
                Console.Error.WriteLine($"random: '{arg}' is not a valid integer");
                Environment.Exit(1);
            }
        }


        if (args.Length == 1)
        {
            // random SEED - set a specific seed and output a random number
            int seed = int.Parse(args[0]);
            var seededRandom = new Random(seed);
            Console.WriteLine(seededRandom.Next(0, 32768));
        }
        else if (args.Length == 2 && int.TryParse(args[0], out int min) && int.TryParse(args[1], out int max))
        {
            // Generate a random integer between min and max (inclusive)
            if (min > max)
            {
                (min, max) = (max, min); // Swap if min > max
            }
            Console.WriteLine(Random.Shared.Next(min, max + 1)); // +1 because Next() is exclusive of upperBound
        }
        else if (args.Length == 3 && int.TryParse(args[0], out int start) &&
                     int.TryParse(args[1], out int step) && int.TryParse(args[2], out int end))
        {
            // random START STEP END - generate a random number from the sequence
            if (step == 0)
            {
                Console.Error.WriteLine("random: step cannot be 0");
                Environment.Exit(1);
            }

            // Ensure start is less than end when step is positive, and start is greater than end when step is negative
            if ((step > 0 && start > end) || (step < 0 && start < end))
            {
                Console.Error.WriteLine($"random: invalid sequence {start} {step} {end}");
                Environment.Exit(1);
            }

            // Calculate how many numbers are in the sequence
            int count = (int)Math.Floor((double)(end - start) / step) + 1;

            // Generate a random index and calculate the corresponding value
            int randomIndex = Random.Shared.Next(0, count);
            int result = start + (step * randomIndex);

            Console.WriteLine(result);
        }
        else
        {
            Console.Error.WriteLine("random: Invalid arguments");
            Usage();
            Environment.Exit(1);
        }
    }

    public static void Usage()
    {
        Console.WriteLine("random - generate random number");
        Console.WriteLine("");
        Console.WriteLine("Usage:");
        Console.WriteLine("  random                 # Random number between 0-32767");
        Console.WriteLine("  random SEED            # Random number with specific seed");
        Console.WriteLine("  random START END       # Random number between START and END (inclusive)");
        Console.WriteLine("  random START STEP END  # Random element from sequence");
        Console.WriteLine("  random choice [ITEMS...]");
        Console.WriteLine("");
        Console.WriteLine("Options:");
        Console.WriteLine("  -h, --help  Display this help message");
        Console.WriteLine("");
        Console.WriteLine("Examples:");
        Console.WriteLine("  random       # Random number between 0 and 32767");
        Console.WriteLine("  random 42    # Random number with seed 42");
        Console.WriteLine("  random 1 10  # Random number between 1 and 10 (inclusive)");
        Console.WriteLine("  random 0 2 10  # Random even number between 0 and 10");
        Console.WriteLine("  random choice apple banana cherry  # Randomly select one item");
    }
}
