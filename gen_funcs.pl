#!/usr/bin/env perl

use strict;
use warnings;
use autodie;
use 5.010;

die "usage $0 [-h] DIR" if @ARGV == 0 or $ARGV[0] =~ /^--?h/;

my $DO_NOT_EDIT_HEADER = '// DO NOT EDIT!!! AUTOMATICALLY GENERATED!!!';

# These functions are variadics, but only with one possible param. In
# this case, discard the variadic property and use a default value for
# this optional parameter.
my %IGNORE_VARIADIC = (Between   => 'BoundsInIn',
		       N         => 0,
		       TruncTime => 0);

my $dir = shift;

opendir(my $dh, $dir);

my %funcs;

while (readdir $dh)
{
    if (/^td_.*\.go\z/ and not /_test.go\z/)
    {
        open(my $fh, '<', "$dir/$_");
        while (defined(my $line = <$fh>))
        {
            if ($line =~ /^func ([A-Z]\w*)\((.*?)\) TestDeep \{$/)
            {
		my $func = $1;
		if ($func ne 'Ignore')
		{
		    my @args;
		    foreach my $arg (split(/, /, $2))
		    {
			my %arg;
			@arg{qw(name type)} = split(/ /, $arg, 2);
			if ($arg{variadic} = $arg{type} =~ s/^\.{3}//)
			{
			    if (exists $IGNORE_VARIADIC{$func})
			    {
				$arg{default} = $IGNORE_VARIADIC{$func};
				delete $arg{variadic};
			    }
			}

			push(@args, \%arg);
		    }
		    $funcs{$func}{args} = \@args;
		}
            }
        }
        close $fh;
    }
}

closedir($dh);

my $funcs_content = <<EOH;
package testdeep

$DO_NOT_EDIT_HEADER

import (
\t"regexp"
\t"testing"
\t"time"
)
EOH

foreach my $func (sort keys %funcs)
{
    my $func_name = "Cmp$func";

    my $cmp_args = 'got interface{}';
    my $call_args = '';

    foreach my $arg (@{$funcs{$func}{args}})
    {
	$call_args .= ', ' unless $call_args eq '';
	$call_args .= $arg->{name};

	$cmp_args .= ", $arg->{name} ";

	if ($arg->{variadic})
	{
	    $call_args .= '...';
	    $cmp_args .= '[]';
	}

	$cmp_args .= $arg->{type};
    }

    $funcs_content .= <<EOF;

// Cmp$func is a shortcut for:
//   CmpDeeply(t, got, $func($call_args), args...)
EOF

    my $last_arg = $funcs{$func}{args}[-1];
    if (exists $last_arg->{default})
    {
	$funcs_content .= <<EOF
//
// $func() optional parameter "$last_arg->{name}" is here mandatory.
// $last_arg->{default} value should be passed to mimic its absence in
// original $func() call.
EOF
    }

    $funcs_content .= <<EOF;
func Cmp$func(t *testing.T, $cmp_args, args ...interface{}) bool {
\treturn CmpDeeply(t, got, $func($call_args), args...)
}
EOF
}

my $examples = do { open(my $efh, '<', 'example_test.go'); local $/; <$efh> };
my $funcs_reg = join('|', keys %funcs);

while ($examples =~ /^func Example($funcs_reg)(_\w+)?\(\) \{\n(.*?)^\}/gms)
{
    push(@{$funcs{$1}{examples}}, { name => $2 // '', code => $3 });
}

open(my $fh, "| gofmt -s > '$dir/funcs.go'");
print $fh $funcs_content;
close $fh;
say "$dir/funcs.go generated";
undef $fh;


my $funcs_test_content = <<EOH;
package testdeep_test

$DO_NOT_EDIT_HEADER

import (
\t"fmt"
\t"regexp"
\t"strconv"
\t"testing"
\t"time"

\t. "github.com/maxatome/go-testdeep"
)
EOH

my($rep, $reb, $rec);
$rep = qr/\( [^()]* (?:(??{ $rep }) [^()]* )* \)/x; # recursively matches (...)
$reb = qr/\[ [^][]* (?:(??{ $reb }) [^][]* )* \]/x; # recursively matches [...]
$rec = qr/\{ [^{}]* (?:(??{ $rec }) [^{}]* )* \}/x; # recursively matches {...}

sub extract_params
{
    my($func, $params_str) = @_;
    my $str = substr($params_str, 1, -1);

    my @params;
    for (;;)
    {
	if ($str =~ /\G\s*
	             ( "(?:\\.|[^"]+)*"            # "string"
	              |`[^`]*`                     # `string`
                      |&[a-zA-Z_]\w*$rec           # &Struct{...}
                      |\[[^][]*\]\w+$rec           # []Array{...}
	              |map${reb}\w+$rec            # map[...]Type{...}
                      |[a-zA-Z_]\w*(?:\.\w+)?(?:$rec|$rep)? # Str{...}, Fn(...), pkg.var
	              |[\w.*+-\/]+                 # 123*pkg.var...
                      )\s*(,|\z)/msgx)
	{
	    push(@params, $1);
	    $2 or return @params;
	}
	else
	{
	    die "Cannot extract params from $func: $params_str\n"
	}
    }
}

foreach my $func (sort keys %funcs)
{
    my $args = $funcs{$func}{args};

    foreach my $example (@{$funcs{$func}{examples}})
    {
	my($name, $code) = @$example{qw(name code)};

        $code =~ s%CmpDeeply\(t,\s+got,\s+$func($rep)%
                   my @params = extract_params("$func$name", $1);
                   my $repl = "Cmp$func(t, got";
                   for (my $i = 0; $i < @$args; $i++)
                   {
                       $repl .= ', ';
                       if ($args->[$i]{variadic})
                       {
                           if (defined $params[$i])
                           {
                               $repl .= '[]' . $args->[$i]{type} . '{'
                                      . join(', ', @params[$i .. $#params])
                                      . '}';
                           }
                           else
                           {
                               $repl .= 'nil';
                           }
                           last
                       }
                       $repl .= $params[$i]
                           // $args->[$i]{default}
			   // die("Internal error, no param: "
				  . "$func$name -> #$i/$args->[$i]{name}!\n");
                   }
                   $repl
                  %egs;

	$funcs_test_content .= <<EOF;

func ExampleCmp$func$name() {
$code}
EOF
    }
}

open($fh, "| gofmt -s > '$dir/funcs_test.go'");
print $fh $funcs_test_content;
close $fh;
say "$dir/funcs_test.go → generated";